package models

import (
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Blog struct {
	// mongodb bson id，类似于主键
	Id      bson.ObjectId
	Email   string
	CDate   time.Time
	Title   string
	Subject string
	//评论数
	CommentCnt int
	//阅读数
	ReadCnt int
	//用于归档
	Year int
}

func (blog *Blog) Validator(v *revel.Validation) {
	v.Check(blog.Title,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{200},
	)
	v.Check(blog.Email,
		revel.Required{},
		revel.MaxSize{50},
	)
	v.Email(blog.Email)
	v.Check(blog.Subject,
		revel.Required{},
		revel.MinSize{1},
	)
}

func (dao *Dao) CreateBlog(blog *Blog) error {
	blogCollection := dao.session.DB(DbName).C(BlogCollection)
	//set the time
	blog.Id = bson.NewObjectId()
	blog.CDate = time.Now()
	blog.Year = blog.CDate.Year()
	_, err := blogCollection.Upsert(bson.M{"_id": blog.Id}, blog)
	if err != nil {
		revel.WARN.Printf("Unable to save blog: %v error %v", blog, err)
	}
	return err
}

func (dao *Dao) FindBlogs() []Blog {
	blogCollection := dao.session.DB(DbName).C(BlogCollection)
	blogs := []Blog{}
	query := blogCollection.Find(bson.M{}).Sort("-cdate").Limit(50)
	query.All(&blogs)
	return blogs
}

func (blog *Blog) GetShortTitle() string {
	if len(blog.Title) > 35 {
		return blog.Title[:35]
	}
	return blog.Title
}

func (blog *Blog) GetShortContent() string {
	if len(blog.Subject) > 200 {
		return blog.Subject[:200]
	}
	return blog.Subject
}
